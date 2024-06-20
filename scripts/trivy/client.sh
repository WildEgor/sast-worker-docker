# $1 - mode (image,config), $2 - server API, $3 - target
/usr/local/bin/trivy $1 --server $2 --quiet --format json $3