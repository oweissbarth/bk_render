"""empty message

Revision ID: 75525bf5d3e5
Revises: 09a5a06b3eef
Create Date: 2016-12-05 10:12:13.671642

"""
from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import mysql

# revision identifiers, used by Alembic.
revision = '75525bf5d3e5'
down_revision = '09a5a06b3eef'
branch_labels = None
depends_on = None


def upgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    op.add_column('task', sa.Column('num_chunks', sa.Integer(), nullable=True))
    op.drop_column('task', 'chunks')
    # ### end Alembic commands ###


def downgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    op.add_column('task', sa.Column('chunks', mysql.INTEGER(display_width=11), autoincrement=False, nullable=True))
    op.drop_column('task', 'num_chunks')
    # ### end Alembic commands ###