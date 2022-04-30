from Cython.Build import cythonize
from setuptools import setup

setup(
    name="wordsmith db",
    ext_modules=cythonize("./wordsmith/db.pyx"),
    zip_safe=False,
)
