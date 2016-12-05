"""empty message

Revision ID: 09a5a06b3eef
Revises: 8f2449a5abe0
Create Date: 2016-12-04 11:39:01.181090

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '09a5a06b3eef'
down_revision = '8f2449a5abe0'
branch_labels = None
depends_on = None


def upgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    op.create_table('worker',
    sa.Column('id', sa.Integer(), nullable=False),
    sa.Column('name', sa.String(length=128), nullable=True),
    sa.Column('lastOnline', sa.DateTime(), nullable=True),
    sa.Column('ip', sa.String(length=64), nullable=True),
    sa.PrimaryKeyConstraint('id')
    )
    op.add_column('chunk', sa.Column('worker', sa.Integer(), nullable=True))
    op.create_foreign_key(None, 'chunk', 'worker', ['worker'], ['id'])
    # ### end Alembic commands ###


def downgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    op.drop_constraint(None, 'chunk', type_='foreignkey')
    op.drop_column('chunk', 'worker')
    op.drop_table('worker')
    # ### end Alembic commands ###