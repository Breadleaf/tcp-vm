programs=("python" "bake" "go" "pkl" "docker")

all_found=true

for program in "${programs[@]}"; do
	if ! which "$program" > /dev/null 2>&1; then
		echo "Error: Could not find '$program'"
		all_found=false
	fi
done

if "$all_found"; then
	echo "Environment is ready"
else
	echo "Environment is missing some dependencies"
fi
