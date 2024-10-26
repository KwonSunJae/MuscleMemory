from setuptools import setup, find_packages

setup(
    name="musclememory",
    version="0.1",
    packages=find_packages(),
    include_package_data=True,
    install_requires=[
        'Click',
        'GitPython',
        'cryptography>=42.0.0',
    ],
    entry_points={
        'console_scripts': [
            'muscle=musclememory.cli:cli',
        ],
    },
)
