mock-gen:
	$(MAKE) -C app/inventory mock-gen
	$(MAKE) -C pkg mock-gen

unit-test:
	$(MAKE) -C app/inventory unit-test
	$(MAKE) -C pkg unit-test

integration-test:
	$(MAKE) -C app/inventory integration-test
	$(MAKE) -C pkg integration-test

goa-gen:
	$(MAKE) -C app/inventory goa-gen-vehicle

swaggo-gen:
	$(MAKE) -C app/inventory swaggo-gen-vehicle