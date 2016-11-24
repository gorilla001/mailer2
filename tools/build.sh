#!/bin/bash
set -e
set -o pipefail

# prepare 
if [ "${IN_CONTAINER}" != "yes" ]; then
	echo;echo "build should run within container"
	exit 1
fi

if [ ! -d "${SRC_DIR}" ]; then
    mkdir -p "${SRC_DIR}"
fi
if [ ! -d "${PRODUCT_DIR}" ]; then
    mkdir -p "${PRODUCT_DIR}"
fi

# link project to GOPATH
mkdir -p $(dirname "/go/src/${PKG}")
ln -sv ${SRC_DIR} /go/src/${PKG}

# golint on the source codes
if [ "${NO_GOLINT}" != "yes" ]; then
	if ! command -v golint; then
		echo;echo "golint not found"
		exit 1
	fi
	# PKG="github.com/tinymailer/mailer/"   # no need, inherit from Build Images Env
	PKGEscape="${PKG//\//\\/}"			# replace all of / to \/
	PKGVendor="$PKG/vendor/"
	PKGVendorEscape="${PKGVendor//\//\\/}"
	packages=( $( 
		go list  -f '{{join .Deps "\n"}}'  $PKG |\
			awk '( $0~/^'$PKGEscape'/ && $0!~/^'$PKGVendorEscape'/ )'
		)
	)
	set +e # turn off
	failNum=0
	for p in ${packages[@]}
	do
		echo; echo " ---> golinting on [$p] ..."
		golint -set_exit_status $p
		if [ $? -ne 0 ]; then
			((failNum++))
		fi
	done
	if [ ${failNum} -gt 0 ]; then
		echo;echo "golint on codes failed"
		exit 1
	fi
	echo;echo "golint passed"
	set -e # turn on

	if [ "${GOLINT_ONLY}" == "yes" ]; then
		exit
	fi
fi

# build the binary

# prepare version info
pushd $SRC_DIR >/dev/null 2>&1
VERSION=$(cat VERSION.txt)
GIT_COMMIT=$(git rev-parse --short HEAD)
if [ -n "$(git status --porcelain --untracked-files=no)" ]; then
	GIT_COMMIT=${GIT_COMMIT}-dirty 
fi
BUILDAT=$(date +%F_%T-%Z)

#go build -a -installsuffix cgo -ldflags="-X main.buildVersion 0.14.0 \
#47             -X main.buildRevision $(git rev-parse --short HEAD) \
#48             -X main.buildBranch master \
#49             -X main.buildUser mountkin@gmail.com \
#50             -X main.buildDate '$(date)' \
#51             -X main.goVersion 1.4.1 -w"  \

# time to build
LD_FLAGS="-X $PKG/version.version=$VERSION -X $PKG/version.gitCommit=$GIT_COMMIT -X $PKG/version.buildAt=$BUILDAT -w"
DstFile="${PRODUCT_DIR}/mailer-${VERSION}-${GIT_COMMIT}"
env CGO_ENABLE=0 GOOS=linux go build -a -ldflags="${LD_FLAGS}" -o $DstFile $PKG

# symbolic link to the latest
pushd ${PRODUCT_DIR} >/dev/null 2>&1
ln -svf ${DstFile##*/} mailer-latest
