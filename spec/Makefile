DRAFT = draft-ramos-aavp-protocol
VERSION = 00

.PHONY: all clean txt html xml

all: txt html

txt: $(DRAFT)-$(VERSION).txt
html: $(DRAFT)-$(VERSION).html
xml: $(DRAFT)-$(VERSION).xml

$(DRAFT)-$(VERSION).xml: $(DRAFT).md
	kramdown-rfc2629 $< > $@

$(DRAFT)-$(VERSION).txt: $(DRAFT)-$(VERSION).xml
	xml2rfc --text $<

$(DRAFT)-$(VERSION).html: $(DRAFT)-$(VERSION).xml
	xml2rfc --html $<

clean:
	rm -f $(DRAFT)-$(VERSION).xml $(DRAFT)-$(VERSION).txt $(DRAFT)-$(VERSION).html
