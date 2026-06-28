import { createFileRoute } from '@tanstack/react-router';

import { useAppForm } from '@/hooks/demo.form';

const emailRegex = /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i;
const zipCodeRegex = /^\d{5}(-\d{4})?$/;
const phoneRegex = /^(\+\d{1,3})?\s?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$/;

const AddressForm = () => {
  const form = useAppForm({
    defaultValues: {
      address: {
        city: '',
        country: '',
        state: '',
        street: '',
        zipCode: '',
      },
      email: '',
      fullName: '',
      phone: '',
    },
    onSubmit: ({ value }) => {
      console.info('Form submitted successfully!', value);
    },
    validators: {
      onBlur: ({ value }) => {
        const errors = {
          fields: {},
        } as {
          fields: Record<string, string>;
        };

        if (value.fullName.trim().length === 0) {
          errors.fields.fullName = 'Full name is required';
        }

        return errors;
      },
    },
  });

  return (
    <main className="demo-page demo-center">
      <section className="demo-panel w-full max-w-2xl">
        <div className="mb-6">
          <p className="island-kicker mb-2">TanStack Form</p>
          <h1 className="demo-title">Address Form</h1>
          <p className="demo-muted mt-2">Nested fields, field-level validation, and a select input.</p>
        </div>
        <form
          className="space-y-6"
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}>
          <form.AppField name="fullName">{(field) => <field.TextField label="Full Name" />}</form.AppField>

          <form.AppField
            name="email"
            validators={{
              onBlur: ({ value }) => {
                if (!value || value.trim().length === 0) {
                  return 'Email is required';
                }
                if (!emailRegex.test(value)) {
                  return 'Invalid email address';
                }

                return;
              },
            }}>
            {(field) => <field.TextField label="Email" />}
          </form.AppField>

          <form.AppField
            name="address.street"
            validators={{
              onBlur: ({ value }) => {
                if (!value || value.trim().length === 0) {
                  return 'Street address is required';
                }

                return;
              },
            }}>
            {(field) => <field.TextField label="Street Address" />}
          </form.AppField>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <form.AppField
              name="address.city"
              validators={{
                onBlur: ({ value }) => {
                  if (!value || value.trim().length === 0) {
                    return 'City is required';
                  }

                  return;
                },
              }}>
              {(field) => <field.TextField label="City" />}
            </form.AppField>
            <form.AppField
              name="address.state"
              validators={{
                onBlur: ({ value }) => {
                  if (!value || value.trim().length === 0) {
                    return 'State is required';
                  }

                  return;
                },
              }}>
              {(field) => <field.TextField label="State" />}
            </form.AppField>
            <form.AppField
              name="address.zipCode"
              validators={{
                onBlur: ({ value }) => {
                  if (!value || value.trim().length === 0) {
                    return 'Zip code is required';
                  }
                  if (!zipCodeRegex.test(value)) {
                    return 'Invalid zip code format';
                  }

                  return;
                },
              }}>
              {(field) => <field.TextField label="Zip Code" />}
            </form.AppField>
          </div>

          <form.AppField
            name="address.country"
            validators={{
              onBlur: ({ value }) => {
                if (!value || value.trim().length === 0) {
                  return 'Country is required';
                }

                return;
              },
            }}>
            {(field) => (
              <field.Select
                label="Country"
                placeholder="Select a country"
                values={[
                  { label: 'United States', value: 'US' },
                  { label: 'Canada', value: 'CA' },
                  { label: 'United Kingdom', value: 'UK' },
                  { label: 'Australia', value: 'AU' },
                  { label: 'Germany', value: 'DE' },
                  { label: 'France', value: 'FR' },
                  { label: 'Japan', value: 'JP' },
                ]}
              />
            )}
          </form.AppField>

          <form.AppField
            name="phone"
            validators={{
              onBlur: ({ value }) => {
                if (!value || value.trim().length === 0) {
                  return 'Phone number is required';
                }
                if (!phoneRegex.test(value)) {
                  return 'Invalid phone number format';
                }

                return;
              },
            }}>
            {(field) => <field.TextField label="Phone" placeholder="123-456-7890" />}
          </form.AppField>

          <div className="flex justify-end">
            <form.AppForm>
              <form.SubscribeButton label="Submit" />
            </form.AppForm>
          </div>
        </form>
      </section>
    </main>
  );
};

export const Route = createFileRoute('/demo/form/address')({
  component: AddressForm,
});
