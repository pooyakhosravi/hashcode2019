# hashcode2019
Google Hashcode competition 2019

This repo contains code for the  Hashcode competition.
In code/ you can find `judge.py` which can take input and output files and score them based on the problem criteria.

A CLI command can be used to either find the score of an input and 1 output:
```
python3 judge.py {input_file} {output_file}
python3 judge.py a.txt a.001.txt
```

Similar structure (with addition of `-f`) can be used to find the max score of many outputs:
```
python3 judge.py -f {input_file} {output_folder}
python3 judge.py -f a.txt out/
```

With the folder mode, verbose can be used to print best score at for each file and the best score so far.
```
python3 judge.py -vf {input_file} {output_folder}
```

