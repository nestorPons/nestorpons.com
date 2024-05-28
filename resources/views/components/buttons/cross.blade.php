@include("components.buttons.template",[
    "icon"=>"xmark",
    "color"=>"danger",
    "class"=>$class??"fs-5 border-0",
    "message"=>$message??null,
    "function"=>$function??null,
    "label"=>$label??"Delete",
    "closeModal"=>$closeModal??false,
    "onclick"=>$onclick??false,
    
    ])