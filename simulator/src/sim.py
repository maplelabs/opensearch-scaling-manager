import sys
import config_parser

if __name__ == '__main__':
    if str(sys.argv[1]) == 'configure':
        config_parser.parse_config(str(sys.argv[2]))
    else:
        sys.stdout.write("Invalid operation")
        exit()