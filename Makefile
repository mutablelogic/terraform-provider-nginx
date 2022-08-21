# Paths to tools needed in dependencies
GO := $(shell which go)

# Paths to locations, etc
BUILD_DIR := "build"
CMD_DIR := $(wildcard cmd/*)
BUILD_FLAGS := 

# Targets
all: clean cmd

cmd: $(CMD_DIR)

test:
	@${GO} mod tidy
	@${GO} test -v ./pkg/...

$(CMD_DIR): dependencies mkdir FORCE
	@echo Build cmd $(notdir $@)
	@${GO} build ${BUILD_FLAGS} -o ${BUILD_DIR}/$(notdir $@) ./$@

FORCE:

dependencies:
ifeq (,${GO})
        $(error "Missing go binary")
endif

mkdir:
	@echo Mkdir ${BUILD_DIR}
	@install -d ${BUILD_DIR}

clean:
	@echo Clean	
	@rm -fr $(BUILD_DIR)
	@${GO} mod tidy
	@${GO} clean

