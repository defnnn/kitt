ENV := env PATH=$(PWD)/../bin:$(PWD)/../exec:$(HOME)/.asdf/bin:$(HOME)/.asdf/installs/bin:$(PATH)

rwildcard = $(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) $(filter $2,$d))

%.stdout: %.test
	$(ENV) ./$< > $@ 2> $(patsubst %.stdout,%.stderr,$@) \
		|| (touch $<; false)
	git diff --exit-code --src-prefix=expected/ --dst-prefix=actual/ \
		$@ $(patsubst %.stdout,%.stderr,$@) \
		|| (touch $<; false)

all:
	touch *.test
	$(MAKE) always

always:
	-$(MAKE) --always-make test

test: $(patsubst %.test,%.stdout,$(call rwildcard,,%.test))

setup:
	git clone -b v0.7.8 https://github.com/asdf-vm/asdf ~/.asdf
	$(ENV) asdf plugin-add letfn--shflags  https://github.com/defn/asdf-plugin
	$(ENV) asdf plugin-add defn--sub https://github.com/defn/asdf-plugin
	$(ENV) asdf install
