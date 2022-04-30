python -m cProfile -o run.pstats run.py
gprof2dot -f pstats run.pstats | dot -Tpng -o run.png