/// Базовый URL API. Для эмулятора Android: 10.0.2.2:8080, для веб — localhost:8080.
String get apiBaseUrl {
  const env = String.fromEnvironment('NC_API_URL', defaultValue: '');
  if (env.isNotEmpty) return env;
  return 'http://localhost:8080';
}

/// Палитра НАРХОЗ (сайт университета)
class NarxozColors {
  static const int dark = 0xFFA82523;
  static const int red = 0xFFD50032;
}
