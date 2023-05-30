# YAML File Schema

The `<field>`s define the structure of each record in the file. Each
row will have some value for every defined `<field>` in the yaml.

Every column field is treated as a string, however the data within that string
can be created by multiple `Faker` functions, each with various `constraints`
to govern the value.

An example configuration can be found in [input/default.yaml](./input/default.yaml).

|Field|Type|Purpose|Constraint|Example|Required|
|-----|----|-------|----------|-------|--------|
|`metadata.rowCount`|int|The total number of rows in the file|`>0`|`5000`| ✅ 
|`metadata.delimiter`|string|The delimiter to use in the file (defaults to `','` when not specified)|Must be surrounded by single quotes|`'\|'`|❌
|`corruption.enabled`|boolean|(NOT IMPLEMENTED) Indicates whether corruption is enabled or not|-|false|❌
|`corruption.chance`|float|The chance that a record will be corrupted|`[0.0, 1.0]`|`0.99` (99% chance)|❌
|`corruption.min`|int|The minimum number of records to corrupt|`[0, <metadata.rowCount>]`|1|❌
|`corruption.max`|int|The maximum number of records to corrupt|`[1, <metadata.rowCount>]`|3|❌
|`corruption.method`|string|The method of corruption. `DATA_TYPE` chooses a different `source` than the specified source for the column. `FORMAT` adds additional delimiters to the row. Required when `corruption.enabled` is `true`.|One of `[DATA_TYPE, FORMAT]`| `DATA_TYPE` |✅ 
|`columns`|list|The list of columns to add to the CSV|-|-|✅
|`columns.<field>.label`|string|The label for the field|Any string with special chars must be surrounded in single quotes|`FIRST_NAME`| ✅
|`columns.<field>.source`|string|The [Faker](https://github.com/brianvoe/gofakeit) method used to generate fake data for the field|Must be a supported source in [faker.go](./app/faker/faker.go)|`PHONE`| ✅
|`columns.<field>.numericalConstraint`|object|Describe numerical constrains for the value of this field|Must define both min and max when included|-|❌
|`columns.<field>.numericalConstraint.min`|int|The minimum value this field can have|-|`1000`|✅
|`columns.<field>.numericalConstraint.max`|int|The maximum value this field can have|-|`9999`|✅
|`columns.<field>.stringConstraint.oneOf`|list|The possible values this field can have|-|`- value1 -value2 -value3`|❌
|`columns.<field>.stringConstraint.regex`|string|The regex pattern used to generate this field|Must be valid regex and surrounded in single quotes|`'[a-zA-Z]{5}'`|❌
|`columns.<field>.timestampConstraint`|object|Describe constraints for a timestamp format and range|Must include format and **BOTH** or **NEITHER** after and before|-|❌
|`columns.<field>.timestampConstraint.format`|string|The date format for the timestamp|Must be parsable by Go's `time.Time` library|`2006-01-02'`|✅
|`columns.<field>.timestampConstraint.after`|string|The beginning of the timestamp range|Must be a valid timestamp|`2023-05-26'`|❌
|`columns.<field>.timestampConstraint.before`|string|The end of the timestamp range|Must be a valid timestamp|`2023-05-27'`|❌
|`columns.<field>.generalConstraint`|object|General constraints that could apply to any data type|Must define `skipChance` when included|-|❌
|`columns.<field>.generalConstraint.skipChance`|float|The chance a column will have an empty value|`[0.0, 1.0]`|`0.01` (1% chance)|✅
