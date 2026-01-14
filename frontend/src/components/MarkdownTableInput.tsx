import React, { useState } from 'react';
import { Input, Button, Space, Alert } from 'antd';
import { useTranslation } from 'react-i18next';

const { TextArea } = Input;

interface MarkdownTableInputProps {
  onSubmit: (data: string) => void;
  placeholder?: string;
  exampleText?: string;
  loading?: boolean;
}

const MarkdownTableInput: React.FC<MarkdownTableInputProps> = ({
  onSubmit,
  placeholder,
  exampleText,
  loading = false,
}) => {
  const { t } = useTranslation();
  const [value, setValue] = useState('');
  const [showExample, setShowExample] = useState(false);

  const handleSubmit = () => {
    if (value.trim()) {
      onSubmit(value);
    }
  };

  const handleClear = () => {
    setValue('');
  };

  return (
    <Space direction="vertical" style={{ width: '100%' }} size="middle">
      {exampleText && (
        <>
          <Button type="link" onClick={() => setShowExample(!showExample)}>
            {showExample ? t('common.hideExample') : t('common.showExample')}
          </Button>
          {showExample && (
            <Alert
              message={t('common.exampleFormat')}
              description={<pre style={{ margin: 0 }}>{exampleText}</pre>}
              type="info"
              closable
              onClose={() => setShowExample(false)}
            />
          )}
        </>
      )}
      <TextArea
        value={value}
        onChange={(e) => setValue(e.target.value)}
        placeholder={placeholder || t('common.enterMarkdownTable')}
        rows={10}
        style={{ fontFamily: 'monospace' }}
      />
      <Space>
        <Button type="primary" onClick={handleSubmit} loading={loading} disabled={!value.trim()}>
          {t('common.submit')}
        </Button>
        <Button onClick={handleClear} disabled={!value.trim()}>
          {t('common.clear')}
        </Button>
      </Space>
    </Space>
  );
};

export default MarkdownTableInput;
