import 'package:flutter/material.dart';

class RegulationsScreen extends StatelessWidget {
  const RegulationsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Регламент и академический календарь', style: Theme.of(context).textTheme.titleLarge),
                const SizedBox(height: 12),
                const Text('Академический календарь и правила обучения размещены на официальном сайте колледжа. Силлабусы и материалы — в разделе «Расписание» (ДЗ и материалы к парам).'),
                const SizedBox(height: 16),
                FilledButton.tonal(
                  onPressed: () {},
                  child: const Text('Сайт колледжа НАРХОЗ'),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
