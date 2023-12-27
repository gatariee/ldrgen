import subprocess
import argparse
import os

def run_ldr(ldr_path, template_path):
    bin_arg = os.path.join(template_path, "bin", "Calc.bin")
    out_arg = "out"
    ldr_arg = "Inline"
    template_arg = template_path
    command = [ldr_path, "-bin", bin_arg, "-out", out_arg, "-ldr", ldr_arg, "--template", template_arg]

    print(f"Running command: {' '.join(command)}")
    result = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, check=True)
    return result.stdout.strip()

def main():
    parser = argparse.ArgumentParser(description="Test Inline LDR")
    parser.add_argument('--ldr_path', type=str, required=True, help="Path to the loader binary")
    parser.add_argument('--template_path', type=str, required=True, help="Path to the template")
    
    args = parser.parse_args()

    output = run_ldr(args.ldr_path, args.template_path)
    expected_substring = 'Done!'
    
    assert expected_substring in output, "Test failed: 'Done!' not found in output"
    print("Test passed: ✅")
    print(f"Loader should be in {args.template_path}out, compile and test it- it should pop calc :)")


if __name__ == '__main__':
    main()
