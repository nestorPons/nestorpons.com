@include("components.buttons.template",[
    "class"=>"px-3 " . ($class ?? ''),
    "collapse"=>$target,
    "icon"=>$icon??"info",
    "color"=>"secondary",
    "message"=>$message??null,
    "function"=>$function??null,
    "label"=>$label??null,
    "closeModal"=>$modal??null
    ])