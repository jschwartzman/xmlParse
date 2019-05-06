#!/usr/bin/env python3
#############################################################################
# file: xmlParse.py
#         A Python SAX utility that parses and displays
#         well-formed xml files.
# author: John Schwartzman, Forte Systems, Inc.
# VERSION_NUMBER = "0.1.0"
# last revision:	03/19/2019
#############################################################################

import xml.sax
import sys

# define some color escape sequences
RED     = '\033[1;31;40m'   #    red on black - for element name
GREEN   = '\033[1;32;40m'   #  green on black - for element data
BLUE    = '\033[1;34;40m'   #   blue on black - for attributes
CYAN    = '\033[1;36;40m'   #   cyan on black - for future use
WHITE   = '\033[0;37;40m'   #  white on black - normal printing
VIOLET  = '\033[1;35;40m'   # violet on black - for future use
BLACK   = '\033[0;30;40m'   #  black on black - for spaces

# define some color macros
SPACES      = BLACK + '%s' + WHITE
STARTNAME   = RED + '%s ' + WHITE
ENDNAME     = RED + '/%s\n' + WHITE
PARENSTMNT  = BLUE + '(%s = %s) ' + WHITE
ELEMENTDATA = GREEN + '%s ' + WHITE

class XmlContentHandler(xml.sax.ContentHandler):
    charBuffer = ''
    elementStack = []
    nLastWritePos = 0

    def __init__(self):
        xml.sax.ContentHandler.__init__(self)

    def pushElementToStack(self, name):      # push name on list
        pos = len(self.elementStack)         # return position before push
        self.elementStack.append(name)
        return pos

    def popElementFromStack(self, name):     # pop name from list
        length = len(self.elementStack)      # return its position after pop
        if name == self.elementStack.pop():
            return length - 1
        else:
            print('%s\nelement %s was not at the top of the stack\n\n' %
                  WHITE, name)
            sys.exit(4)

    def writeSpaces(self, pos, chars):      # position the cursor column
        for _ in range(0, pos): # don't need index so use dummy variable
            sys.stdout.write(SPACES % chars)

    def writeNewLine(self):                 # position the cursor row
        print('')

    def writeStartName(self, name):         # write name at current x,y pos
        sys.stdout.write(STARTNAME % name)

    def writeEndName(self, name):           # write name at current x,y pos
        sys.stdout.write(ENDNAME % name)

    def writeElementData(self, data):       # write content at current x,y pos
        if data != "":
            sys.stdout.write(ELEMENTDATA % data)

    def writeAttributes(self, attributes):  # write element attributes at x,y
        for iname in attributes.getNames():
            sys.stdout.write(PARENSTMNT % (iname, attributes.getValue(iname)))

    def getCharacterData(self):             # get the charBuffer contents
        data = self.charBuffer.strip()      # and clear charBuffer
        self.charBuffer = ''
        return data

    def startElement(self, name, attributes):   # we've encountered a startElement
        pos = self.pushElementToStack(name)
        self.writeNewLine()
        self.writeSpaces(pos, '   ')            # write 3 spaces per index pos
        self.writeStartName(name)
        self.writeAttributes(attributes)
        self.nLastWritePos = pos

    def endElement(self, name):                 # we've encountered an endElement
        pos = self.popElementFromStack(name)
        charStr = self.getCharacterData()
        self.writeElementData(charStr)  # write it
        if pos < self.nLastWritePos:   # write name at curent x pos?
            self.writeSpaces(pos, '   ')    # position endElement on x-axis
        self.writeEndName(name)

    def characters(self, content):              # push characters into charBuffer
        self.charBuffer += content

def usage(exitFlag):
    print ('\nUSAGE: ./xmlParse.py xmlFileToView\n\n')
    sys.exit(exitFlag)

def main(sourceFileName):
    try:
        source = open(sourceFileName)
        # instantiate and initialize XmlContentHandler
        xml.sax.parse(source, XmlContentHandler())
    except xml.sax.SAXParseException as e:  # handle parser error
        print('\nFailed to parse ' + sourceFileName +
            '. This appears not to be a valid xml document: ' +
            e.getMessage() + '\n\n')
        sys.exit(2)
    except OSError as e:    # handle os error
        print('\nFailed to open ' + sourceFileName + ': ' + str(e) + '\n\n')
        sys.exit(3)

if __name__ == "__main__":
    if sys.argv[1] == "-h" or sys.argv[1] == "--help":  # does user want help
        usage(0)
    if len(sys.argv) != 2:  # there must be 2 arguments
        usage(1)

    main(sys.argv[1])
    print('')
    sys.exit(0)
