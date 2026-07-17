# Justfile — Common tasks for ash development

# Variables
VERSION := $(shell date +%Y.%m.%d)

# Default target
default: help

help:
	@echo "ash — Arch Snapshot Hypervisor"
	@echo ""
	@echo "Targets:"
	@echo "  build-iso       Build ISO locally (requires root, Arch host)"
	@echo "  test-iso        Test ISO boot in QEMU"
	@echo "  build-vm        Build all VM formats via Packer"
	@echo "  build-cloud     Build cloud images (AWS/GCP/Azure)"
	@echo "  sign            Generate signatures + SLSA provenance"
	@echo "  distribute      Upload to mirrors (R2, Bunny, Archive.org, Torrents)"
	@echo "  release         Full release pipeline (build + sign + distribute)"
	@echo "  landing         Build and preview landing page"
	@echo "  clean           Clean build artifacts"
	@echo ""

# Build ISO (requires root on Arch Linux)
build-iso:
	sudo ./scripts/build-iso.sh $(VERSION)

# Test ISO boot
test-iso:
	./scripts/test-iso.sh out/ash-$(VERSION).iso

# Build VM formats
build-vm:
	cd packer && packer init . && packer build -var "version=$(VERSION)" -var "iso_path=../out/ash-$(VERSION).iso" ash-iso.pkr.hcl

# Build cloud images
build-cloud:
	cd packer && for f in aws-ami.pkr.hcl gcp-image.pkr.hcl azure-image.pkr.hcl; do \
	  [ -f "$$f" ] && packer build -var "version=$(VERSION)" -var "iso_path=../out/ash-$(VERSION).iso" "$$f" || true; \
	done

# Generate signatures and SLSA provenance
sign:
	./scripts/sign-provenance.sh $(VERSION)

# Distribute to mirrors
distribute:
	./scripts/distribute.sh $(VERSION)

# Full release pipeline
release: build-iso test-iso sign build-vm build-cloud distribute
	@echo "Release $(VERSION) complete!"

# Landing page
landing:
	cd landing-page && npm install && npm run dev

landing-build:
	cd landing-page && npm install && npm run build

landing-preview:
	cd landing-page && npm run preview

# Clean
clean:
	rm -rf out/
	rm -rf packer/output-*/
	rm -rf landing-page/dist/ landing-page/.astro/
	docker system prune -f

.PHONY: help build-iso test-iso build-vm build-cloud sign distribute release landing landing-build landing-preview clean