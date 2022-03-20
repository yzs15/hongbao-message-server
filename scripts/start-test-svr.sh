sed -i '' "s/[\#]\{0,\}. scripts\/conf-2C4C.sh/\. scripts\/conf-2C4C.sh/g" scripts/test-all.sh
sed -i '' "s/[\#]\{0,\}. scripts\/conf-noise.sh/#. scripts\/conf-noise.sh/g" scripts/test-all.sh

export FLASK_APP=test-server
flask run -h 0.0.0.0 -p 9002