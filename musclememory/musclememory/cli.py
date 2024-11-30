import click
from .commands.init import init
from .commands.ready import ready
from .commands.add import add
from .commands.enroll import enroll
from .commands.complete import complete
from .commands.config import config_group
from .commands.fetch import fetch
from .commands.history import history
from .commands.rollback import rollback

@click.group()
def cli():
    """MuscleMemory CLI"""
    pass

# 명령어 등록
cli.add_command(init)
cli.add_command(ready)
cli.add_command(add)
cli.add_command(enroll)
cli.add_command(complete)
cli.add_command(config_group)
cli.add_command(fetch)
cli.add_command(history)
cli.add_command(rollback)

if __name__ == "__main__":
    cli()
