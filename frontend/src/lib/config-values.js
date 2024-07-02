export const ValueTypes = {
	STRING: 0,
	INTEGER: 1,
	FLOAT: 2,
	BOOLEAN: 3
};

export const getValue = (configValue) => {
		switch (configValue.ValueType) {
			case ValueTypes.STRING:
				return configValue.StrValue ?? '';
			case ValueTypes.INTEGER:
				return configValue.IntValue ?? 0;
			case ValueTypes.FLOAT:
				return configValue.FloatValue ?? 0.0;
			case ValueTypes.BOOLEAN:
				return configValue.BoolValue ?? false;
			default:
			console.error(configValue);
				throw new Error('Somehow reached unreachable code!');
		}
}

export const updateValue = (configValue, newValue) => {
		switch (configValue.ValueType) {
			case ValueTypes.STRING:
				configValue.StrValue = newValue;
				break;
			case ValueTypes.INTEGER:
				configValue.IntValue = newValue;
				break;
			case ValueTypes.FLOAT:
				configValue.FloatValue = newValue;
				break;
			case ValueTypes.BOOLEAN:
				configValue.BoolValue = newValue;
				break;
			default:
			console.error(configValue);
				throw new Error(`Somehow reached unreachable code!`);
		}
}
