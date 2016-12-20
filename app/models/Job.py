from app import db
from flask import url_for, abort


class Job(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    name = db.Column(db.String(64))
    startFrame = db.Column(db.Integer)
    endFrame = db.Column(db.Integer)
    num_chunks = db.Column(db.Integer)
    filename = db.Column(db.String(128))
    chunks = db.relationship('Chunk', backref='job', lazy='dynamic', cascade="delete")

    def __init__(self):
        pass

    def from_json(self, json):
        try:
            self.name = json['name']
            self.startFrame = json['startFrame']
            self.endFrame = (json['endFrame'])
            self.num_chunks = (json['numChunks'])

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
        return url_for("getJob", id=self.id, _external=True)

    def __repr__(self):
        return '<Job %r>' % self.name
