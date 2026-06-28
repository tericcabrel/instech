import { createFileRoute } from '@tanstack/react-router';
import { z } from 'zod';

import { useAppForm } from '@/hooks/demo.form';

const schema = z.object({
  description: z.string().min(1, 'Description is required'),
  title: z.string().min(1, 'Title is required'),
});

const SimpleForm = () => {
  const form = useAppForm({
    defaultValues: {
      description: '',
      title: '',
    },
    onSubmit: ({ value }) => {
      console.info('Form submitted successfully!', value);
    },
    validators: {
      onBlur: schema,
    },
  });

  return (
    <main className="demo-page demo-center">
      <section className="demo-panel w-full max-w-2xl">
        <div className="mb-6">
          <p className="island-kicker mb-2">TanStack Form</p>
          <h1 className="demo-title">Simple Form</h1>
          <p className="demo-muted mt-2">A small validated form using the generated field helpers.</p>
        </div>
        <form
          className="space-y-6"
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}>
          <form.AppField name="title">{(field) => <field.TextField label="Title" />}</form.AppField>

          <form.AppField name="description">{(field) => <field.TextArea label="Description" />}</form.AppField>

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

export const Route = createFileRoute('/demo/form/simple')({
  component: SimpleForm,
});
