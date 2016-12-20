"""empty message

Revision ID: e8beb0f08b0d
Revises: d4fdd86f11fe
Create Date: 2016-12-20 10:13:38.929414

"""
from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import mysql

# revision identifiers, used by Alembic.
revision = 'e8beb0f08b0d'
down_revision = 'd4fdd86f11fe'
branch_labels = None
depends_on = None


def upgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    op.add_column('job', sa.Column('filename', sa.String(length=128), nullable=True))
    op.drop_column('job', 'filePath')
    # ### end Alembic commands ###


def downgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    op.add_column('job', sa.Column('filePath', mysql.VARCHAR(length=128), nullable=True))
    op.drop_column('job', 'filename')
    # ### end Alembic commands ###