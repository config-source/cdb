export const ValueType = {
	STRING: 0,
	INTEGER: 1,
	FLOAT: 2,
	BOOLEAN: 3
};

export const getValue = (configValue) => {
	switch (configValue.ValueType) {
		case ValueType.STRING:
			return configValue.StrValue ?? '';
		case ValueType.INTEGER:
			return configValue.IntValue ?? 0;
		case ValueType.FLOAT:
			return configValue.FloatValue ?? 0.0;
		case ValueType.BOOLEAN:
			return configValue.BoolValue ?? false;
		default:
			console.error(configValue);
			throw new Error('Somehow reached unreachable code!');
	}
};

export const updateValue = (configValue, newRawValue) => {
	// TODO: type guard newValue
	switch (configValue.ValueType) {
		case ValueType.STRING:
			configValue.StrValue = newRawValue;
			configValue.IntValue = null;
			configValue.FloatValue = null;
			configValue.BoolValue = null;
			break;
		case ValueType.INTEGER:
			configValue.IntValue = newRawValue;
			configValue.StrValue = null;
			configValue.FloatValue = null;
			configValue.BoolValue = null;
			break;
		case ValueType.FLOAT:
			configValue.FloatValue = newRawValue;
			configValue.StrValue = null;
			configValue.IntValue = null;
			configValue.BoolValue = null;
			break;
		case ValueType.BOOLEAN:
			configValue.BoolValue = newRawValue;
			configValue.StrValue = null;
			configValue.IntValue = null;
			configValue.FloatValue = null;
			break;
		default:
			console.error(configValue);
			throw new Error(`Somehow reached unreachable code!`);
	}
};
