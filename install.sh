RELEASES_URL="https://api.github.com/repos/holy-filipp/sshsync/releases"
UNAME=$(uname)

if [ "$UNAME" == "Linux" ] ; then
  DOWNLOAD_PATH=$(curl $RELEASES_URL | jq -r '.[0].assets[] | select(.name | endswith("linux-amd64.tar.gz")) | .browser_download_url')
elif [ "$UNAME" == "Darwin" ] ; then
	DOWNLOAD_PATH=$(curl $RELEASES_URL | jq -r '.[0].assets[] | select(.name | endswith("darwin-arm64.tar.gz")) | .browser_download_url')
fi

wget -qO- "$DOWNLOAD_PATH" | gunzip | tar xvf -