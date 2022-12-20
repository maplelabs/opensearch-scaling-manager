import sys
import os
from pathlib import Path
from unittest import expectedFailure
import yaml
import pytest


current_file = Path(__file__).parent.parent.resolve()
test_path = str(os.path.join(str(current_file), "tests", "config_test"))
src_path = os.path.join(str(current_file), "src")
sys.path.insert(0, src_path)

from config_parser import validate_config, parse_config, Config
from errors import ValidationError
from constants import CONFIG_FILE_PATH
import constants as const


"""Checks the config file with all the required datas"""
def test_validate_config():
    with open(os.path.join(test_path, "config_1_P.yaml"), "r") as file:
        is_valid, errors = validate_config(yaml.safe_load(file))
        assert is_valid == True
        assert errors == {}


"""Validates if config has searches field in it."""
def test_validate_config_without_searches():
    with open(os.path.join(test_path, "config_2_P.yaml"), "r") as file:
        is_valid, errors = validate_config(yaml.safe_load(file))
        assert is_valid == True
        assert errors == {}


"""Validates if config has searches field in it."""
def test_validate_config_without_data_ingestion():
    with open(os.path.join(test_path, "config_3_P.yaml"), "r") as file:
        is_valid, errors = validate_config(yaml.safe_load(file))
        assert is_valid == True
        assert errors == {}


"""Validates if config has searches field in it""" 
def test_validate_config_missing_parameter():
    with open(os.path.join(test_path, "config_4_F.yaml"), "r") as file:
        is_valid, errors = validate_config(yaml.safe_load(file))
        assert is_valid == False
        assert errors == {
            "index_clean_up_age_days": ["required field"],
            "index_roll_over_size": ["required field"],
        }


"""Checks if the config has a valid data type"""
def test_validate_config_invalid_data_type():
    with open(os.path.join(test_path, "config_5_F.yaml"), "r") as file:
        is_valid, errors = validate_config(yaml.safe_load(file))
        assert is_valid == False
        assert errors == {
            "data_ingestion": [
                {
                    "states": [
                        {3: [{"ingestion_rate_gb_per_hr": ["must be of number type"]}]}
                    ]
                }
            ]
        }


"""Validates config against the list of dictionary in schema"""
def test_validate_config_missing_nested_key():
    with open(os.path.join(test_path, "config_6_F.yaml"), "r") as file:
        is_valid, errors = validate_config(yaml.safe_load(file))
        assert is_valid == False
        assert errors == {"searches": [{0: [{"probability": ["required field"]}]}]}


"""Checks if it's a valid config and place the file in the simulator/src/main and return"""
def test_parse_config():
    fp = open(os.path.join(test_path, "config_1_P.yaml"), "r")
    all_configs = yaml.safe_load(fp.read())

    data_generation_interval_minutes = all_configs.pop(
        const.DATA_GENERATION_INTERVAL_MINUTES
    )
    data_ingestion = all_configs.pop(const.DATA_INGESTION)
    searches = all_configs.pop(const.SEARCHES)
    stats = all_configs

    expected_config = Config(
        stats, data_ingestion, searches, data_generation_interval_minutes
    )
    config = parse_config(os.path.join(test_path, "config_1_P.yaml"))
    assert (
        expected_config.data_generation_interval_minutes
        == config.data_generation_interval_minutes
    )
    assert expected_config.stats == config.stats
    assert expected_config.data_ingestion == config.data_ingestion
    assert expected_config.searches == config.searches


"""Checks the config is complete or not """
def test_parse_config_error_reading_config():
    with pytest.raises(ValidationError) as e:
        parse_config(os.path.join(test_path, "config_5_F.yaml"))
        assert "error reading config file - " == e


"""If required field is not there in config and dont place in src path"""
def test_parse_config_validate_error():
    with pytest.raises(ValidationError) as e:
        parse_config(os.path.join(test_path, "config_4_F.yaml"))
        assert "Error validating config file - " == e
