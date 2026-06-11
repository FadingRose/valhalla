function token
    set -l value $argv[1]
    set -l decimals $argv[2]

    # Ensure the string is long enough by padding with leading zeros
    # e.g. value "10", decimals 3 -> becomes "0010" to produce "0.010"
    set -l min_len (math $decimals + 1)
    set value (string pad -w $min_len -c 0 -- $value)

    set -l len (string length -- $value)
    set -l cut_point (math $len - $decimals)

    # Split into integer and fractional parts
    set -l int_part (string sub -l $cut_point -- $value)
    set -l dec_part (string sub -s (math $cut_point + 1) -- $value)

    echo "$int_part.$dec_part"
end
