all: test

# Run tests.
test:
	go test -v ./...

# Push the git tag.
git-push-tag: VERSION=$(shell cat VERSION)
git-push-tag:
	git push origin ${VERSION}

# Add the tag.
git-tag-release: VERSION=$(shell cat VERSION)
git-tag-release: check-release-version
	git tag --annotate ${VERSION} --message "go-certcentral ${VERSION}"

# Check whether the tag exists already.
check-release-version: VERSION=$(shell cat VERSION)
check-release-version:
	if test x$$(git tag --list ${VERSION}) != x; \
	then \
		echo "Tag [${VERSION}] already exists. Please check the working copy."; git diff . ; exit 1;\
	fi

# Tag a new release of the library.
release: VERSION=$(shell cat VERSION)
release: git-tag-release git-push-tag

# Clean up
clean:
	rm -rf vendor

.PHONY: vendor
vendor:
	go mod vendor
