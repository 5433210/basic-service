#/bin/bash

# db2struct
base_dir=./
host=localhost
package=model
dbname=authn
dbuser=root
dbpassword=rootroot
target=target

get_all_tables_sql="select table_name from information_schema.tables where table_schema='$dbname'"

tables=`mysql -u$dbuser -p$dbpassword -e "$get_all_tables_sql" -s -N 2>/dev/null`

output_dir=$base_dir/$dbname
mkdir -p $output_dir

table_array=($(echo $tables | tr " " "\n"))
for table_name in "${table_array[@]}"
do
entity_name=`echo $table_name | gsed -e 's/^tb//' -e 's/\(_\)\(.\)/\2/g'`
struct_name=`echo $table_name | gsed -e 's/^tb//' -e 's/\(_\)\(.\)/\U\2/g'`
target=$output_dir/$entity_name.go
db2struct --host=$host --database=$dbname --table=$table_name --package=$package --struct=$struct_name --password=$dbpassword --user=$dbuser --json --gorm --target=$target
goimports -w $target
gofmt -w $target
done

commandInPath()
{
    echo $get_all_tables_sql
    echo $tables
}

# commandInPath

