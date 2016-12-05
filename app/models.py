from flask.ext.sqlalchemy import SQLAlchemy
from datetime import datetime
from flask import url_for, abort, g, send_from_directory
from itsdangerous import (TimedJSONWebSignatureSerializer
                          as Serializer, BadSignature, SignatureExpired)

from app import db

from app import app


class Task(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(64))
    startFrame = db.Column(db.Integer)
    endFrame = db.Column(db.Integer)
    num_chunks = db.Column(db.Integer)
    filePath = db.Column(db.String(128))
    chunks = db.relationship('Chunk', backref='task', lazy='dynamic')

    def __init__(self):
        pass

    def from_json(self, json):
        try:
            self.name = json['name']
            self.startFrame = json['startFrame']
            self.endFrame = (json['endFrame'])

        except KeyError:
            print("Invalid format")
            abort(400)

        return self

    def to_json(self):
        return{
            "id": self.id,
            "name": self.name,
            "startFrame": self.startFrame,
            "endFrame": self.endFrame
        }

    def get_url(self):
        return url_for("getTask", id=self.id, _external=True)

    def __repr__(self):
        return '<Task %r>' % self.name


class Chunk(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    startFrame = db.Column(db.Integer)
    endFrame = db.Column(db.Integer)
    taskId = db.Column(db.Integer, db.ForeignKey("task.id"))
    available = db.Column(db.Boolean)
    done = db.Column(db.Boolean)
    worker = db.Column(db.Integer, db.ForeignKey("worker.id"))
    jobFile = db.Column(db.String(128))

    def __init__(self):
        pass

    def __init__(self, t, start, end):
        self.taskId = t.id
        self.startFrame = start
        self.endFrame = end
        self.available = True
        self.jobFile = url_for("getJobfile", id=t.id, _external=True)

    def from_json(self, json):
        try:
            self.name = json['name']
            self.startFrame = json['startFrame']
            self.endFrame = json['endFrame']

        except KeyError:
            print("Invalid format")
            abort(400)

        return self

    def to_json(self):
        return{
            "id": self.id,
            "startFrame": self.startFrame,
            "endFrame": self.endFrame,
            "jobFile": self.jobFile
        }

    def get_url(self):
        return url_for("getChunk", id=self.id, _external=True)

    def __repr__(self):
        return '<Chunk %r>' % self.id


class Worker(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(128))
    lastOnline = db.Column(db.DateTime)
    ip = db.Column(db.String(64))

    def __init__(self):
        self.lastOnline = datetime.utcnow()

    def from_json(self, json):
        try:
            self.name = json['name']
            self.ip = (json['ip'])

        except KeyError:
            print("Invalid format")
            abort(400)

        return self

    def to_json(self):
        return{
            "id": self.id,
            "name": self.name,
            "lastOnline": self.lastOnline
        }

    def get_url(self):
        return url_for("getWorker", id=self.id, _external=True)

    def __repr__(self):
        return '<Worker %r>' % self.id
