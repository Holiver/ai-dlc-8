// Storage service for managing localStorage and sessionStorage

const storageService = {
  // LocalStorage methods
  local: {
    set: (key: string, value: any): void => {
      try {
        const serializedValue = JSON.stringify(value);
        localStorage.setItem(key, serializedValue);
      } catch (error) {
        console.error('Error saving to localStorage:', error);
      }
    },

    get: <T = any>(key: string): T | null => {
      try {
        const item = localStorage.getItem(key);
        return item ? JSON.parse(item) : null;
      } catch (error) {
        console.error('Error reading from localStorage:', error);
        return null;
      }
    },

    remove: (key: string): void => {
      try {
        localStorage.removeItem(key);
      } catch (error) {
        console.error('Error removing from localStorage:', error);
      }
    },

    clear: (): void => {
      try {
        localStorage.clear();
      } catch (error) {
        console.error('Error clearing localStorage:', error);
      }
    },
  },

  // SessionStorage methods
  session: {
    set: (key: string, value: any): void => {
      try {
        const serializedValue = JSON.stringify(value);
        sessionStorage.setItem(key, serializedValue);
      } catch (error) {
        console.error('Error saving to sessionStorage:', error);
      }
    },

    get: <T = any>(key: string): T | null => {
      try {
        const item = sessionStorage.getItem(key);
        return item ? JSON.parse(item) : null;
      } catch (error) {
        console.error('Error reading from sessionStorage:', error);
        return null;
      }
    },

    remove: (key: string): void => {
      try {
        sessionStorage.removeItem(key);
      } catch (error) {
        console.error('Error removing from sessionStorage:', error);
      }
    },

    clear: (): void => {
      try {
        sessionStorage.clear();
      } catch (error) {
        console.error('Error clearing sessionStorage:', error);
      }
    },
  },
};

export default storageService;
