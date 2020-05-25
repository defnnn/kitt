ENV := env PATH=$(PWD)/../bin:$(PWD)/../exec:$(PATH)

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

test: $(patsubst %.test,%.stdout,$(call rwildcard,,%.test))

always:
	$(MAKE) --always-make test
