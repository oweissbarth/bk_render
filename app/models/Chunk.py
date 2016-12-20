from flask.ext.sqlalchemy import SQLAlchemy
from datetime import datetime
from flask import url_for, abort, g, send_from_directory
from itsdangerous import (TimedJSONWebSignatureSerializer
                          as Serializer, BadSignature, SignatureExpired)

import os

from app import db

from app import app

class Chunk(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    startFrame = db.Column(db.Integer)
    endFrame = db.Column(db.Integer)
    jobId = db.Column(db.Integer, db.ForeignKey("job.id"))
    available = db.Column(db.Boolean)
    done = db.Column(db.Boolean)
    worker = db.Column(db.Integer, db.ForeignKey("worker.id"))
    jobFile = db.Column(db.String(128))
    jobFileType = db.Column(db.String(20))

    def __init__(self):
        pass

    def __init__(self, t, start, end):
        self.jobId = t.id
        self.startFrame = start
        self.endFrame = end
        self.available = True
        self.jobFile = url_for("getJobfile", id=t.id, _external=True)
        filename, ext = os.path.splitext(t.filename)
        self.jobFileType = ext

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
            "jobId": self.jobId,
            "startFrame": self.startFrame,
            "endFrame": self.endFrame,
            "jobFile": self.jobFile,
            "jobFileType": self.jobFileType
        }

    def get_url(self):
        return url_for("getChunk", id=self.id, _external=True)

    def __repr__(self):
        return '<Chunk %r>' % self.id
