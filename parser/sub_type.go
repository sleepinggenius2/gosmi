package parser

/*
<integerSubType>
::= <empty>
  | "(" <range> ["|" <range>]... ")"

<octetStringSubType>
::= <empty>
  | "(" "SIZE" "(" <range> ["|" <range>]... ")" ")"

<range>
::= <value>
  | <value> ".." <value>

<value>
::= "-" <number>
  | <number>
  | <hexString>
  | <binString>

where:
<empty>     is the empty string
<number>    is a non-negative integer
<hexString> is a hexadecimal string (e.g., '0F0F'H)
<binString> is a binary string (e.g, '1010'B)

<range> is further restricted as follows:
	- any <value> used in a SIZE clause must be non-negative.
	- when a pair of values is specified, the first value
	  must be less than the second value.
	- when multiple ranges are specified, the ranges may
	  not overlap but may touch. For example, (1..4 | 4..9)
	  is invalid, and (1..4 | 5..9) is valid.
	- the ranges must be a subset of the maximum range of the
	  base type.

	  Some examples of legal sub-typing:

	  Integer32 (-20..100)
	  Integer32 (0..100 | 300..500)
	  Integer32 (300..500 | 0..100)
	  Integer32 (0 | 2 | 4 | 6 | 8 | 10)
	  OCTET STRING (SIZE(0..100))
	  OCTET STRING (SIZE(0..100 | 300..500))
	  OCTET STRING (SIZE(0 | 2 | 4 | 6 | 8 | 10))
	  SYNTAX   TimeInterval (0..100)
	  SYNTAX   DisplayString (SIZE(0..32))

(Note the last two examples above are not valid in a TEXTUAL
CONVENTION, see [3].)

Some examples of illegal sub-typing:

  Integer32 (150..100)         -- first greater than second
  Integer32 (0..100 | 50..500) -- ranges overlap
  Integer32 (0 | 2 | 0 )       -- value duplicated
  Integer32 (MIN..-1 | 1..MAX) -- MIN and MAX not allowed
  Integer32 (SIZE (0..34))     -- must not use SIZE
  OCTET STRING (0..100)        -- must use SIZE
  OCTET STRING (SIZE(-10..100)) -- negative SIZE
*/