i=0
while ((i++ < 20)); do
  mkdir -p "hellopy-$i"
  cp hellopy/__init__.py "hellopy-$i"/__init__.py
  cp hellopy/function.json "hellopy-$i"/function.json
  cp hellopy/sample.dat "hellopy-$i"/sample.dat
done
