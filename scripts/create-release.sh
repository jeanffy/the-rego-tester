if [ -z "$1" ]; then
  echo "No version supplied"
  exit 1
fi

version=$1

echo "Creating tag for version $version"

git tag -d $version
git push --delete origin $version
git tag $version -a -m "Version $version"
git push origin $version
