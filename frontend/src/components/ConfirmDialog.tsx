import { Modal } from 'antd';
import { ExclamationCircleOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';

interface ConfirmDialogOptions {
  title: string;
  content: string;
  onOk: () => void | Promise<void>;
  onCancel?: () => void;
  okText?: string;
  cancelText?: string;
  type?: 'info' | 'success' | 'error' | 'warning' | 'confirm';
}

export const useConfirmDialog = () => {
  const { t } = useTranslation();

  const showConfirm = ({
    title,
    content,
    onOk,
    onCancel,
    okText,
    cancelText,
    type = 'confirm',
  }: ConfirmDialogOptions) => {
    Modal.confirm({
      title,
      content,
      icon: <ExclamationCircleOutlined />,
      okText: okText || t('common.confirm'),
      cancelText: cancelText || t('common.cancel'),
      onOk,
      onCancel,
      okButtonProps: {
        danger: type === 'error',
      },
    });
  };

  return { showConfirm };
};

export default useConfirmDialog;
