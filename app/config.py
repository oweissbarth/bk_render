import os

basedir = os.path.abspath(os.path.dirname(__file__))

SQLALCHEMY_DATABASE_URI = 'mysql://bc:bc@localhost/bc'
SQLALCHEMY_MIGRATE_REPO = os.path.join(basedir, 'db_repository')
SECRET_KEY = '.B\xe9\xb6\xf0\x1d\xb4\xb3\x9eRAHj\x1e\xbfK\r\x95\xe5;@\x82\xba\xbd'

DEBUG = True

UPLOAD_FOLDER = "/home/oliver/uploads/"
