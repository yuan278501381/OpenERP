/**
 * Unified Logger for Frontend Shell
 */

type LogLevel = 'info' | 'warn' | 'error' | 'debug';

class Logger {
  private generateTraceId(): string {
    return Math.random().toString(36).substring(2, 15);
  }

  private log(level: LogLevel, message: string, data?: any) {
    const traceId = this.generateTraceId();
    const timestamp = new Date().toISOString();
    const logPrefix = `[${timestamp}] [${level.toUpperCase()}] [TraceID: ${traceId}]`;
    
    switch (level) {
      case 'info':
        console.info(logPrefix, message, data ? data : '');
        break;
      case 'warn':
        console.warn(logPrefix, message, data ? data : '');
        break;
      case 'error':
        console.error(logPrefix, message, data ? data : '');
        break;
      case 'debug':
        console.debug(logPrefix, message, data ? data : '');
        break;
    }
  }

  info(message: string, data?: any) {
    this.log('info', message, data);
  }

  warn(message: string, data?: any) {
    this.log('warn', message, data);
  }

  error(message: string, data?: any) {
    this.log('error', message, data);
  }

  debug(message: string, data?: any) {
    this.log('debug', message, data);
  }
}

export const logger = new Logger();
