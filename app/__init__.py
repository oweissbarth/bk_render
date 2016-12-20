from flask import Flask
from flask import jsonify, json, Response
from flask import abort
from flask import render_template
from flask import request
from flask import g
from werkzeug.utils import secure_filename

import os

from flask.ext.sqlalchemy import SQLAlchemy


app = Flask(__name__)

app.config.from_pyfile('config.py')

db = SQLAlchemy(app)


from app.models import *


@app.route("/")
def index():
    jobs = Job.query.all()
    workers = Worker.query.all()
    return render_template("index.html", jobs=jobs, workers=workers)


@app.route("/jobs", methods=["GET"])
def getJobs():
    return to_json(Job.query.all())


@app.route("/jobs/<int:id>", methods=["GET"])
def getJob(id):
    return to_json(Job.query.filter_by(id=id).first())


@app.route("/jobs/<int:id>", methods=["DELETE"])
def deleteJob(id):
    job = Job.query.filter_by(id=id).first()

    if(job is None):
        abort(404)

    db.session.delete(job)
    db.session.commit()

    filename = job.filename

    os.remove(os.path.join(app.config['UPLOAD_FOLDER'], filename))

    response = jsonify({})
    response.status_code = 201

    return response


@app.route("/jobs", methods=["POST"])
def addJob():
    json = request.get_json(force=True, silent=True)
    try:
        name = json["name"]
        startFrame = json["startFrame"]
        endFrame = json["endFrame"]
        num_chunks = json["numChunks"]
    except:
        print("invalid format")
        abort(400)

    job = Job().from_json(json)

    db.session.add(job)
    db.session.commit()

    print(job.get_url())
    response = to_json(job.get_url())

    response.status_code = 201

    return response


@app.route("/jobs/<int:id>", methods=["POST"])
def addFile(id):
    t = Job.query.filter_by(id=id).first()

    if(t is None):
        abort(404)

    f = request.files['file']

    extension = f.filename.rsplit('.', 1)[1]
    filename = str(t.id) + "." + extension
    filepath = os.path.join(app.config['UPLOAD_FOLDER'], filename)
    f.save(filepath)

    t.filename = filename

    num_chunks = t.num_chunks
    start = t.startFrame
    end = t.endFrame

    numFrames = end - start
    for i in range(num_chunks):
        chunk = Chunk(t, start + i * int(numFrames / num_chunks),
                      start + (i + 1) * int(numFrames / num_chunks))
        db.session.add(chunk)

    db.session.commit()
    response = jsonify({})
    response.status_code = 200

    return response


@app.route("/jobfile/<int:id>", methods=["GET"])
def getJobfile(id):
    t = Job.query.filter_by(id=id).first()
    return send_from_directory(app.config['UPLOAD_FOLDER'], t.filename,
                               as_attachment=True,
                               attachment_filename="job.blend")


@app.route("/worker/<int:id>/job", methods=["GET"])
def requestJob(id):
    worker = Worker.query.filter_by(id=id).first()
    worker.lastOnline = datetime.utcnow()
    db.session.commit()

    # if there is an unfinished job assigned to the worker he will get it again
    chunk = Chunk.query.filter_by(done=False, worker=id).first()
    if(chunk is None):
        chunk = Chunk.query.filter_by(available=True).first()

    if(chunk is None):
        abort(404)

    chunk.available = False
    chunk.worker = id
    db.session.commit()
    response = to_json(chunk.to_json())
    response.status_code = 200
    return response


@app.route("/worker/<int:workerid>/job/<int:jobid>", methods=["PUT"])
def updateJob(workerid, jobid):
    chunk = Chunk.query.filter_by(id=jobid).first()
    worker = Worker.query.filter_by(id=workerid).first()
    worker.lastOnline = datetime.utcnow()
    db.session.commit()

    if(chunk is None):
        abort(404)
    if(chunk.worker != workerid):
        abort(403)

    json = request.get_json(force=True, silent=True)
    try:
        done = json["done"]
    except:
        print("invalid format")
        abort(400)

    chunk.available = False
    chunk.done = done
    db.session.commit()
    response = to_json(chunk.to_json())
    response.status_code = 200
    return response


@app.route("/worker", methods=["POST"])
def addWorker():
    json = request.get_json(force=True)
    try:
        name = json["name"]
        ip = json["ip"]
    except:
        print("invalid format")
        abort(400)

    worker = Worker.query.filter_by(name=name).first()
    if(worker is None):
        worker = Worker().from_json(json)
        db.session.add(worker)
        db.session.commit()

    response = jsonify(worker.to_json())
    response.status_code = 201

    return response


def to_json(input):
    return Response(json.dumps(input), mimetype='application/json')
