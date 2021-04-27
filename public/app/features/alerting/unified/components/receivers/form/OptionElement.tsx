import React, { FC } from 'react';
import { Checkbox, Field, Input, InputControl, Select, TextArea } from '@grafana/ui';
import { NotificationChannelOption } from 'app/types';
import { useFormContext, FieldError, NestDataObject } from 'react-hook-form';

interface Props {
  option: NotificationChannelOption;
  invalid?: boolean;
  pathPrefix?: string;
  error?: FieldError | NestDataObject<any, FieldError>;
}

export const OptionElement: FC<Props> = ({ option, invalid, pathPrefix = '', error }) => {
  return (
    <Field
      label={option.element !== 'checkbox' ? option.label : undefined}
      description={option.description || undefined}
      invalid={!!error}
      error={error?.message}
    >
      <OptionInput option={option} invalid={invalid} pathPrefix={pathPrefix} />
    </Field>
  );
};

const OptionInput: FC<Props> = ({ option, invalid, pathPrefix = '' }) => {
  const { control, register } = useFormContext();
  const name = `${pathPrefix}${option.propertyName}`;
  switch (option.element) {
    case 'checkbox':
      return <Checkbox name={name} ref={register()} label={option.label} description={option.description} />;
    case 'input':
      return (
        <Input
          invalid={invalid}
          type={option.inputType}
          name={name}
          ref={register({
            required: option.required ? 'Required' : false,
            validate: (v) => (option.validationRule !== '' ? validateOption(v, option.validationRule) : true),
          })}
          placeholder={option.placeholder}
        />
      );

    case 'select':
      return (
        <InputControl
          as={Select}
          options={option.selectOptions}
          control={control}
          name={name}
          invalid={invalid}
          onChange={(values) => values[0].value}
        />
      );

    case 'textarea':
      return (
        <TextArea
          invalid={invalid}
          name={name}
          ref={register({
            required: option.required ? 'Required' : false,
            validate: (v) => (option.validationRule !== '' ? validateOption(v, option.validationRule) : true),
          })}
        />
      );

    default:
      console.error('Element not supported', option.element);
      return null;
  }
};

const validateOption = (value: string, validationRule: string) => {
  return RegExp(validationRule).test(value) ? true : 'Invalid format';
};
