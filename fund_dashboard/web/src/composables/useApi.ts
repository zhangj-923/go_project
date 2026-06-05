export function useApi() {
  const request = async <T>(url: string, options?: RequestInit): Promise<T> => {
    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
      });
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const data = await response.json();
      return data as T;
    } catch (error) {
      console.error(`API Error on ${url}:`, error);
      throw error;
    }
  };

  const get = <T>(url: string) => request<T>(url);
  
  const post = <T>(url: string, body: any) => 
    request<T>(url, { method: 'POST', body: JSON.stringify(body) });

  const del = <T>(url: string) => 
    request<T>(url, { method: 'DELETE' });

  return { get, post, del };
}
