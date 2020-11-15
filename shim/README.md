This module is responsible for standing between dependencies and E code to
allow for relatively easy switching out of dependencies should there ever
be a need, and also to clean up E code by making dependencies' APIs fit more.
Any dependency can be replaced simply by replacing the file standing between
it and E as long as the new file provides the same functions.