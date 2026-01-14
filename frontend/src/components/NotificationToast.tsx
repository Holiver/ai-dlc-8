import { notification } from 'antd';
import type { ArgsProps } from 'antd/es/notification';

type NotificationType = 'success' | 'info' | 'warning' | 'error';

interface NotificationOptions extends Omit<ArgsProps, 'type'> {
  type?: NotificationType;
}

class NotificationService {
  private static instance: NotificationService;

  private constructor() {
    // Configure default notification settings
    notification.config({
      placement: 'topRight',
      duration: 3,
    });
  }

  public static getInstance(): NotificationService {
    if (!NotificationService.instance) {
      NotificationService.instance = new NotificationService();
    }
    return NotificationService.instance;
  }

  public show({ type = 'info', message, description, ...rest }: NotificationOptions) {
    notification[type]({
      message,
      description,
      ...rest,
    });
  }

  public success(message: string, description?: string) {
    this.show({ type: 'success', message, description });
  }

  public error(message: string, description?: string) {
    this.show({ type: 'error', message, description });
  }

  public warning(message: string, description?: string) {
    this.show({ type: 'warning', message, description });
  }

  public info(message: string, description?: string) {
    this.show({ type: 'info', message, description });
  }
}

export const notificationService = NotificationService.getInstance();

export default notificationService;
