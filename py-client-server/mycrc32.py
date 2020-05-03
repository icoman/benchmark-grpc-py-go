import sys



# from golang package crc32


# IEEE is by far and away the most common CRC-32 polynomial.
# Used by ethernet (IEEE 802.3), v.42, fddi, gzip, zip, png, ...
IEEE = 0xedb88320

# Castagnoli's polynomial, used in iSCSI.
# Has better error detection characteristics than IEEE.
# https://dx.doi.org/10.1109/26.231911
Castagnoli = 0x82f63b78

# Koopman's polynomial.
# Also has better error detection characteristics than IEEE.
# https://dx.doi.org/10.1109/DSN.2002.1028931
Koopman = 0xeb31d82e

#crctable = _GenerateCRCTable(IEEE)
crctable = None



def _GenerateCRCTable(poly):
    #
    # from lib/Python27/Lib/zipfile.py
    #
    """
        Generate a CRC-32 table.
        ZIP encryption uses the CRC32 one-byte primitive for scrambling some
        internal keys. We noticed that a direct implementation is faster than
        relying on binascii.crc32().
    """
    table = [0] * 256
    for i in range(256):
        crc = i
        for _ in range(8):
            if crc & 1:
                crc = ((crc >> 1) & 0x7FFFFFFF) ^ poly
            else:
                crc = ((crc >> 1) & 0x7FFFFFFF)
        table[i] = crc
    return table 


def crc32(data):
    global crctable
    if crctable is None:
        crctable = _GenerateCRCTable(IEEE)
    accum = 0
    accum = ~accum
    if sys.version_info >= (3, 0):
        for b in data:
            accum = crctable[(accum ^ b) & 0xFF] ^ ((accum >> 8) & 0x00FFFFFF)
    else:
        for b in data:
            accum = crctable[(accum ^ ord(b)) & 0xFF] ^ ((accum >> 8) & 0x00FFFFFF)
    accum = ~accum
    return accum & 0xFFFFFFFF

