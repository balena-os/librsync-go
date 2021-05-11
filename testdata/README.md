# librsync-go Test Data

Reference files (`*.signature`, `*.delta`) were created using the original (C
version) `rdiff`.

Old and new files created as follows:

* `000.old`/`000.new`: Both files are equal.
* `001.old`/`001.new`: The new file was created by appending some data to the
  old file.
* `002.old`/`002.new`: The new file was created by prepending some data to the
  old file.
* `003.old`/`003.new`: The new file was created by inserting some data in the
  middle of the old file.
* `004.old`/`004.new`: Files of same size, with some smallish sequences of bytes
  arbitrarily changed on the new one.
* `005.old`/`005.new`: New file was created by removing some data from the
  beginning, middle and end of the old file.
* `006.old`/`006.new`: Tiny files crafted to exercise the case in which there
  is a match of the final block (with length less than the block length). This
  happens when using a block length of 2.
* `007.old`/`007.new`: Tiny files crafted to test the case in which the block
  length is larger than the new file. This happens when using a block length of
  5 or more.
* `008.old`/`008.new`: Old file has data, new file is empty.
* `009.old`/`009.new`: Old file is empty, new file has data.
* `010.old`/`010.new`: Both files are empty.
