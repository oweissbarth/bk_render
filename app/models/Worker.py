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

    def online(self):
        self.lastOnline = datetime.utcnow()

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
