#!/usr/bin/env python

from flask_script import Manager
from flask_migrate import Migrate, MigrateCommand
from flask_sqlalchemy import SQLAlchemy
from sqlalchemy.exc import IntegrityError


from app import app
from app import db
from app import models


manager = Manager(app)

migrate = Migrate(app, db)

manager = Manager(app)
manager.add_command("db", MigrateCommand)


@manager.command
def runserver():
    """Start a debug server"""
    app.run(debug=True)


if __name__ == "__main__":
    manager.run()
