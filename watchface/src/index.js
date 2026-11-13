const digits = [
  'd0.png', 'd1.png', 'd2.png', 'd3.png', 'd4.png',
  'd5.png', 'd6.png', 'd7.png', 'd8.png', 'd9.png'
];
const small = [
  'n0.png', 'n1.png', 'n2.png', 'n3.png', 'n4.png',
  'n5.png', 'n6.png', 'n7.png', 'n8.png', 'n9.png'
];
const weekdays = [
  'w0.png', 'w1.png', 'w2.png', 'w3.png',
  'w4.png', 'w5.png', 'w6.png'
];

WatchFace({
  build() {
    const ui = hmUI;
    const widget = ui.widget;
    const normal = ui.show_level.ONLY_NORMAL;
    const left = ui.align.LEFT;
    const image = (src, x, y) => ui.createWidget(widget.IMG, {
      src,
      x,
      y,
      show_level: normal
    });

    ui.createWidget(widget.FILL_RECT, {
      x: 0,
      y: 0,
      w: 390,
      h: 450,
      color: 0,
      show_level: normal
    });

    image('s0.png', 12, 102);
    image('s1.png', 72, 42);
    image('s2.png', 72, 354);
    image('s3.png', 324, 102);

    // 後に作ったウィジェットが前面に描かれる。文字盤内の文字を針で隠さないよう、針を先に作る。
    ui.createWidget(widget.TIME_POINTER, {
      hour_centerX: 195,
      hour_centerY: 225,
      hour_posX: 6,
      hour_posY: 102,
      hour_path: 'h.png',
      minute_centerX: 195,
      minute_centerY: 225,
      minute_posX: 8,
      minute_posY: 166,
      minute_path: 'm.png',
      second_centerX: 195,
      second_centerY: 225,
      second_posX: 6,
      second_posY: 178,
      second_path: 's.png',
      show_level: normal
    });

    ui.createWidget(widget.IMG_TIME, {
      hour_startX: 131,
      hour_startY: 113,
      hour_array: digits,
      hour_zero: 1,
      hour_space: 0,
      hour_follow: 0,
      hour_align: left,
      hour_unit_en: 'colon.png',
      hour_unit_sc: 'colon.png',
      hour_unit_tc: 'colon.png',
      minute_startX: 204,
      minute_startY: 113,
      minute_array: digits,
      minute_zero: 1,
      minute_space: 0,
      minute_follow: 0,
      minute_align: left,
      show_level: normal
    });

    ui.createWidget(widget.IMG_DATE, {
      month_startX: 149,
      month_startY: 288,
      month_zero: 1,
      month_space: 0,
      month_align: left,
      month_follow: 0,
      month_unit_en: 'slash.png',
      month_unit_sc: 'slash.png',
      month_unit_tc: 'slash.png',
      month_sc_array: small,
      month_tc_array: small,
      month_en_array: small,
      day_startX: 187,
      day_startY: 288,
      day_zero: 1,
      day_space: 0,
      day_align: left,
      day_follow: 0,
      day_sc_array: small,
      day_tc_array: small,
      day_en_array: small,
      show_level: normal
    });

    ui.createWidget(widget.IMG_WEEK, {
      x: 220,
      y: 286,
      week_en: weekdays,
      week_sc: weekdays,
      week_tc: weekdays,
      show_level: normal
    });

    ui.createWidget(widget.TEXT_IMG, {
      x: 145,
      y: 316,
      w: 100,
      h: 22,
      type: ui.data_type.BATTERY,
      font_array: small,
      h_space: 0,
      align_h: ui.align.CENTER_H,
      unit_en: 'u.png',
      unit_sc: 'u.png',
      unit_tc: 'u.png',
      show_level: normal
    });
  }
});
