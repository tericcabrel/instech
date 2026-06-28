import { useStore } from '@tanstack/react-form';

import { useFieldContext, useFormContext } from '@/hooks/demo.form-context';

const ErrorMessages = ({ errors }: { errors: Array<string | { message: string }> }) => (
  <>
    {errors.map((error) => (
      <div className="mt-1 text-sm font-semibold text-red-600" key={typeof error === 'string' ? error : error.message}>
        {typeof error === 'string' ? error : error.message}
      </div>
    ))}
  </>
);

export const SubscribeButton = ({ label }: { label: string }) => {
  const form = useFormContext();

  return (
    <form.Subscribe selector={(state) => state.isSubmitting}>
      {(isSubmitting) => (
        <button className="demo-button" disabled={isSubmitting} type="submit">
          {label}
        </button>
      )}
    </form.Subscribe>
  );
};

export const TextField = ({ label, placeholder }: { label: string; placeholder?: string }) => {
  const field = useFieldContext<string>();
  const errors = useStore(field.store, (state) => state.meta.errors);

  return (
    <div>
      <label className="mb-2 block text-sm font-semibold text-[var(--sea-ink)]" htmlFor={label}>
        {label}
        <input
          className="demo-input mt-2"
          onBlur={field.handleBlur}
          onChange={(e) => field.handleChange(e.target.value)}
          placeholder={placeholder}
          value={field.state.value}
        />
      </label>
      {field.state.meta.isTouched && <ErrorMessages errors={errors} />}
    </div>
  );
};

export const TextArea = ({ label, rows = 3 }: { label: string; rows?: number }) => {
  const field = useFieldContext<string>();
  const errors = useStore(field.store, (state) => state.meta.errors);

  return (
    <div>
      <label className="mb-2 block text-sm font-semibold text-[var(--sea-ink)]" htmlFor={label}>
        {label}
        <textarea
          className="demo-textarea mt-2"
          onBlur={field.handleBlur}
          onChange={(e) => field.handleChange(e.target.value)}
          rows={rows}
          value={field.state.value}
        />
      </label>
      {field.state.meta.isTouched && <ErrorMessages errors={errors} />}
    </div>
  );
};

export const Select = ({
  label,
  values,
}: {
  label: string;
  placeholder?: string;
  values: Array<{ label: string; value: string }>;
}) => {
  const field = useFieldContext<string>();
  const errors = useStore(field.store, (state) => state.meta.errors);

  return (
    <div>
      <label className="mb-2 block text-sm font-semibold text-[var(--sea-ink)]" htmlFor={label}>
        {label}
      </label>
      <select
        className="demo-select"
        name={field.name}
        onBlur={field.handleBlur}
        onChange={(e) => field.handleChange(e.target.value)}
        value={field.state.value}>
        {values.map((value) => (
          <option key={value.value} value={value.value}>
            {value.label}
          </option>
        ))}
      </select>
      {field.state.meta.isTouched && <ErrorMessages errors={errors} />}
    </div>
  );
};
