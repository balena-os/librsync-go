# generate test data

```bash
bash ./create-pm-data.sh
```

# perform measurements

```bash
bash ./run-pm.sh
```

# generate graphs

Install dependencies:
```bash
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

```bash
./visualize.py sigTime old.csv new.csv
```
